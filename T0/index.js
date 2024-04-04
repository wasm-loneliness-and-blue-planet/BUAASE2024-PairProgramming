import fs from "fs";
import "./wasm_exec.js";
const wasmBuffer = fs.readFileSync("./main.wasm");
const go = new Go();
const module = await WebAssembly.instantiate(wasmBuffer, go.importObject);
go.run(module.instance);
const { add, getState, setState } = module.instance.exports;
// const { splitter, splitter2, arrGetter, objGetter } = globalThis;
console.log(add(1, 2));
console.log(getState());
setState(114514);
console.log(getState());

function perf(func, ...args) {
  const start = performance.now();
  for (let i = 0; i < 1000; i++) func(...args);
  const end = performance.now();
  return end - start;
}

console.log(splitter("12,345, 6789,12,124,44,55"));
console.log(splitter2("12,345, 6789,12,124,44,55", 7));
console.log(arrGetter([13, 14, 15, 21, 11, 16], 5));
console.log(
  objGetter({ a: 1, b: [true, 666, "hahaha", { key: "value" }], c: 3 }, "b")
);
