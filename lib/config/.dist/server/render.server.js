// lib/config/.dist/server/render.server.js
var __defProp = Object.defineProperty;
var __defNormalProp = (obj, key, value) => key in obj ? __defProp(obj, key, { enumerable: true, configurable: true, writable: true, value }) : obj[key] = value;
var __publicField = (obj, key, value) => __defNormalProp(obj, typeof key !== "symbol" ? key + "" : key, value);
var HYDRATION_START = "[";
var HYDRATION_END = "]";
var CONTENT_REGEX = /[&<]/g;
function escape_html(value, is_attr) {
  const str = String(value ?? "");
  const pattern = CONTENT_REGEX;
  pattern.lastIndex = 0;
  let escaped = "";
  let last = 0;
  while (pattern.test(str)) {
    const i = pattern.lastIndex - 1;
    const ch = str[i];
    escaped += str.substring(last, i) + (ch === "&" ? "&amp;" : ch === '"' ? "&quot;" : "&lt;");
    last = i + 1;
  }
  return escaped + str.substring(last);
}
function lifecycle_outside_component(name) {
  {
    throw new Error(`https://svelte.dev/e/lifecycle_outside_component`);
  }
}
var current_component = null;
function getContext(key) {
  const context_map = get_or_init_context_map();
  const result = (
    /** @type {T} */
    context_map.get(key)
  );
  return result;
}
function setContext(key, context) {
  get_or_init_context_map().set(key, context);
  return context;
}
function get_or_init_context_map(name) {
  if (current_component === null) {
    lifecycle_outside_component();
  }
  return current_component.c ?? (current_component.c = new Map(get_parent_context(current_component) || void 0));
}
function push(fn) {
  current_component = { p: current_component, c: null, d: null };
}
function pop() {
  var component = (
    /** @type {Component} */
    current_component
  );
  var ondestroy = component.d;
  if (ondestroy) {
    on_destroy.push(...ondestroy);
  }
  current_component = component.p;
}
function get_parent_context(component_context) {
  let parent = component_context.p;
  while (parent !== null) {
    const context_map = parent.c;
    if (context_map !== null) {
      return context_map;
    }
    parent = parent.p;
  }
  return null;
}
var BLOCK_OPEN = `<!--${HYDRATION_START}-->`;
var BLOCK_CLOSE = `<!--${HYDRATION_END}-->`;
var HeadPayload = class {
  constructor(css = /* @__PURE__ */ new Set(), out = "", title = "", uid = () => "") {
    __publicField(this, "css", /* @__PURE__ */ new Set());
    __publicField(this, "out", "");
    __publicField(this, "uid", () => "");
    __publicField(this, "title", "");
    this.css = css;
    this.out = out;
    this.title = title;
    this.uid = uid;
  }
};
var Payload = class {
  constructor(id_prefix = "") {
    __publicField(this, "css", /* @__PURE__ */ new Set());
    __publicField(this, "out", "");
    __publicField(this, "uid", () => "");
    __publicField(this, "head", new HeadPayload());
    this.uid = props_id_generator(id_prefix);
    this.head.uid = this.uid;
  }
};
function props_id_generator(prefix) {
  let uid = 1;
  return () => `${prefix}s${uid++}`;
}
var on_destroy = [];
function render$1(component, options = {}) {
  const payload = new Payload(options.idPrefix ? options.idPrefix + "-" : "");
  const prev_on_destroy = on_destroy;
  on_destroy = [];
  payload.out += BLOCK_OPEN;
  if (options.context) {
    push();
    current_component.c = options.context;
  }
  component(payload, options.props ?? {}, {}, {});
  if (options.context) {
    pop();
  }
  payload.out += BLOCK_CLOSE;
  for (const cleanup of on_destroy) cleanup();
  on_destroy = prev_on_destroy;
  let head2 = payload.head.out + payload.head.title;
  for (const { hash, code } of payload.css) {
    head2 += `<style id="${hash}">${code}</style>`;
  }
  return {
    head: head2,
    html: payload.out,
    body: payload.out
  };
}
function head(payload, fn) {
  const head_payload = payload.head;
  head_payload.out += BLOCK_OPEN;
  fn(head_payload);
  head_payload.out += BLOCK_CLOSE;
}
var $$css = {
  hash: "svelte-17r7htw",
  code: '.content.svelte-17r7htw {position:fixed;left:0;right:0;top:0;bottom:0;background:#1e1e2e;color:cadetblue;display:grid;justify-content:center;align-content:center;font-family:"Noto Sans Gothic",serif;}'
};
function Layout($$payload, $$props) {
  $$payload.css.add($$css);
  const { children } = $$props;
  head($$payload, ($$payload2) => {
    $$payload2.out += `<meta charset="UTF-8"/> <meta name="viewport" content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0"/>`;
  });
  $$payload.out += `<div class="content svelte-17r7htw">`;
  children($$payload);
  $$payload.out += `<!----></div>`;
}
function Router($$payload, $$props) {
  push();
  getContext("server");
  pop();
}
function View($$payload, $$props) {
  push();
  const server = getContext("server");
  head($$payload, ($$payload2) => {
    $$payload2.title = `<title>Welcome</title>`;
  });
  Router();
  $$payload.out += `<!----> `;
  Layout($$payload, {
    children: ($$payload2) => {
      $$payload2.out += `<h1>Hello ${escape_html(server.data.name)}.</h1>`;
    }
  });
  $$payload.out += `<!---->`;
  pop();
}
function Render_server($$payload, $$props) {
  push();
  let { id: idInput, data: dataInput, ids: idsInput } = $$props;
  const server = { id: idInput, data: dataInput, ids: idsInput };
  setContext("server", server);
  if ("welcome" === server.id) {
    $$payload.out += "<!--[-->";
    View($$payload);
  } else {
    $$payload.out += "<!--[!-->";
  }
  $$payload.out += `<!--]-->`;
  pop();
}
async function render(props) {
  return render$1(Render_server, { props });
}
export {
  render
};
