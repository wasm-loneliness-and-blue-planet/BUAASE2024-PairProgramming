import fs from "fs";
import "./wasm_exec.js";
const wasmBuffer = fs.readFileSync("./main.wasm");
const go = new Go();
const module = await WebAssembly.instantiate(wasmBuffer, go.importObject);
go.run(module.instance)
const { add, getState, setState } = module.instance.exports
console.log(add(1, 2))
console.log(getState())
setState(114514)
console.log(getState())
