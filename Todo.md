## Phase 1: Environment & CLI Setup

Before you write any proxy logic, you need to handle the "CLI" part of the CLI tool.

- [x] **Initialize your project:** Create your folder and run `npm init` (Node) or set up your virtual env (Python).
- [x] **Argument Parsing:** Set up the logic to capture `--port`, `--origin`, and `--clear-cache`.
- _Goal:_ If you run `tool --port 3000`, your code should be able to print `3000` to the console.

- [x] **The "Exit Early" Command:** Implement the `--clear-cache` logic first.
- _Logic:_ If this flag is present, delete `cache.json` or overwrite it with `{}` and then `process.exit()`.

## Phase 2: The Storage Engine (JSON)

Since you aren't using Redis, you need a way to talk to your file safely.

- [ ] **Create helper functions:**
- `getCache()`: Reads the JSON file and returns an object. Returns `{}` if the file doesn't exist.
- `updateCache(key, data)`: Adds a new response to the object and saves it back to the file.

- [ ] **Define your "Key":** Decide what makes a request unique.
- _Hint:_ Usually, it's just the URL path (e.g., `/products`).

## Phase 3: The Proxy Server

This is the heart of the tool. You need to start a server that listens on the user-provided `--port`.

- [ ] **Start the HTTP Server:** Create a basic server that logs "Request received" for any URL.
- [ ] **Implement the "HIT" Logic:** \* Check your JSON object for the incoming URL.
- If found: Return the stored body and set the header `X-Cache: HIT`.

- [ ] **Implement the "MISS" Logic:**
- If not found: Use a fetch library (like `axios` or `node-fetch`) to call `--origin` + `path`.
- _Example:_ If origin is `http://dummyjson.com` and path is `/users`, fetch `http://dummyjson.com/users`.

- [ ] **Save & Respond:** \* Store that new data in your JSON file.
- Return the data to the user with the header `X-Cache: MISS`.

## Phase 4: Polish & Headers

The requirements specifically ask for the original headers to be preserved.

- [ ] **Header Forwarding:** When you fetch from the origin, make sure to copy the `Content-Type` (like `application/json`) from the origin's response and send it back to your user.
- [ ] **Error Handling:** What happens if the origin server is down? Add a `try/catch` block to return a 500 error instead of crashing your CLI.

## Phase 5: Testing

Run these manual tests to ensure everything works:

- [ ] **Test 1:** Run the tool, make a request, check if `cache.json` was created.
- [ ] **Test 2:** Make the same request again. It should be much faster and show `X-Cache: HIT`.
- [ ] **Test 3:** Run the `--clear-cache` command. Verify `cache.json` is empty.
- [ ] **Test 4:** Make the request again. It should show `X-Cache: MISS` again.

---

### Pro-Tip for the JSON File

When saving the data to your JSON file, use `JSON.stringify(data, null, 2)`. The `null, 2` part will "pretty-print" the JSON, making it much easier for you to read and debug while you're building!

**Which programming language are you planning to use? I can give you the specific libraries that will make the CLI part easiest.**
