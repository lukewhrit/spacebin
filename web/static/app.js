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
