declare global {
  var arrGetter: (arr: Array<T>, index: number) => T;
  var objGetter: (obj: Record<string, T>, key: string) => T;
}