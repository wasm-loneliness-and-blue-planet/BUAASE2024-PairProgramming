import assert from "assert";

// Choose proper "import" depending on your PL.
// import { mancalaResult } from "./t2-as/build/release.js";
// import { mancala_result as mancalaResult } from "./t2_rust/pkg/t1_rust.js"
// [Write your own "import" for other PLs.]

import fs from "fs";
import "../wasm_exec.js";
const wasmBuffer = fs.readFileSync("./main.wasm");
const go = new Go();
const module = await WebAssembly.instantiate(wasmBuffer, go.importObject);
go.run(module.instance);
const { mancalaResult } = globalThis;

assert.strictEqual(mancalaResult(1,[11,12],2),30001);
assert.strictEqual(mancalaResult(1,[14],1),20001);

console.log("🎉 You have passed all the tests provided.");
