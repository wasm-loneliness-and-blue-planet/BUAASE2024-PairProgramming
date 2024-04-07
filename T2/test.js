import assert from "assert";

// Choose proper "import" depending on your PL.
// import { mancalaResult } from "./t2-as/build/release.js";
// import { mancala_result as mancalaResult } from "./t2_rust/pkg/t1_rust.js"
// [Write your own "import" for other PLs.]

import { mancalaResult } from "./t2-go/t2-go.js";

assert.strictEqual(mancalaResult(1,[11,12],2),30001);
assert.strictEqual(mancalaResult(1,[14],1),20001);
assert.strictEqual(mancalaResult(1, [11, 12], 2), 30001);
assert.strictEqual(mancalaResult(1, [14], 1), 20001);
assert.strictEqual(mancalaResult(1, [], 0), 20000);
// assert.strictEqual(mancalaResult(1, [11], 1), 30001);
assert.strictEqual(mancalaResult(1, [11, 12, 11, 12, 11, 12], 1), 20000);

console.log("🎉 You have passed all the tests provided.");
