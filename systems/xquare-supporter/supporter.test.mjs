import assert from "node:assert/strict";
import test from "node:test";

import { run, systemName } from "./supporter.mjs";

test("run writes the system name", () => {
  const output = [];

  run((value) => output.push(value));

  assert.deepEqual(output, [systemName]);
});
