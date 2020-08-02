declare module 'has-value' {
  interface Options {
    isValid: (key: string, obj: object) => boolean;
    split: (splitter: { [Symbol.split](string: string, limit?: number): string[] }, limit?: number) => string[];
    separator: string | RegExp;
    join: (separator?: string) => string;
    joinChar: string;
  }

  export default function (obj: object, path: string, options?: Options): boolean
}
