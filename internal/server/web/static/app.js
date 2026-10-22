// Fixes tab key in textarea
document.querySelector('textarea')?.addEventListener('keydown', function (e) {
  if (e.key.toLowerCase() === 'tab') {
    e.preventDefault();

    const start = this.selectionStart;
    const end = this.selectionEnd;

    // Set textarea value to: text before caret + tab + text after caret
    this.value = this.value.substring(0, start) + '\t' + this.value.substring(end);

    // Move caret to right position
    this.selectionStart = this.selectionEnd = start + 1;
  }
});

// Allows for saving with CTRL+S and CMD+S
document.addEventListener('keydown', function(e) {
  const isMac = navigator.platform.toUpperCase().indexOf('MAC') >= 0;
  if (e.key.toLowerCase() === 's' && (isMac ? e.metaKey : e.ctrlKey)) {
    e.preventDefault();

    // requestSubmit, unlike submit, runs the submit handlers below, so encryption can't be skipped
    document.querySelector('#text')?.requestSubmit();
  }
});

function switchFont(to) {
  const main = document.querySelector('.wysiwyg');

  if (to === 'sans') {
    main.classList.remove('font-serif', 'font-sans');
    main.classList.add('font-sans');

    document.querySelector('#serif').classList.remove('active');
    document.querySelector('#sans').classList.add('active');
  } else if (to === 'serif') {
    main.classList.remove('font-serif', 'font-sans');
    main.classList.add('font-serif');

    document.querySelector('#sans').classList.remove('active');
    document.querySelector('#serif').classList.add('active');
  }
}

// Encrypting new documents. Shown only where the browser can encrypt, and the form is sent as normal
// unless a password has been chosen.
(function () {
  const form = document.querySelector('#text');
  const toggle = document.querySelector('#encrypt');
  const password = document.querySelector('#password');

  if (!form || !toggle || !password || !spacebinCrypto.supported()) {
    return;
  }

  toggle.hidden = false;

  toggle.addEventListener('click', function () {
    const enabled = password.hidden;

    password.hidden = !enabled;
    toggle.setAttribute('aria-pressed', String(enabled));

    if (enabled) {
      password.focus();
    } else {
      password.value = '';
    }
  });

  password.addEventListener('keydown', function (e) {
    if (e.key === 'Enter') {
      e.preventDefault();
      form.requestSubmit();
    }
  });

  form.addEventListener('submit', async function (e) {
    if (password.hidden) {
      return;
    }

    // Never fall back to sending the plaintext
    e.preventDefault();

    if (password.value === '') {
      alert('Enter a password to encrypt this document, or click the lock to turn encryption off.');
      password.focus();
      return;
    }

    try {
      const content = await spacebinCrypto.encrypt(form.elements.content.value, password.value);
      const response = await fetch('/api/', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ content, encrypted: true }),
      });
      const body = await response.json();

      if (body.error) {
        throw new Error(body.error);
      }

      window.location.href = '/' + body.payload.id;
    } catch (err) {
      alert('Could not save the encrypted document: ' + err.message);
    }
  });
})();

// Decrypting encrypted documents
(function () {
  const main = document.querySelector('#encrypted');

  if (!main) {
    return;
  }

  const form = document.querySelector('#decrypt-form');
  const password = document.querySelector('#decrypt-password');
  const error = document.querySelector('#decrypt-error');
  const output = document.querySelector('#decrypted');
  const copy = document.querySelector('#copy');

  if (!spacebinCrypto.supported()) {
    error.textContent = 'This browser can only decrypt documents over HTTPS.';
    error.hidden = false;
    form.hidden = false;
    form.querySelector('[type=submit]').disabled = true;
    return;
  }

  form.hidden = false;

  form.addEventListener('submit', async function (e) {
    e.preventDefault();
    error.hidden = true;

    let plaintext;

    try {
      plaintext = await spacebinCrypto.decrypt(main.dataset.envelope, password.value);
    } catch {
      error.textContent = 'Wrong password.';
      error.hidden = false;
      password.select();
      return;
    }

    // textContent, never innerHTML: the plaintext is untrusted
    output.querySelector('code').textContent = plaintext;
    output.hidden = false;
    document.querySelector('#form').hidden = true;

    copy.hidden = false;
    copy.addEventListener('click', function () {
      navigator.clipboard.writeText(plaintext);
    });
  });
})();
