const go = new Go();
WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
  go.run(result.instance).then(r => r);
  const queryS = window.location.search.split("?")[1];
  if (!queryS) {
    document.innerHTML = JSON.stringify({status: 400, data: "No arguments"});
  } else {
    const query = parseQuery(queryS);
    if (typeof query.code === "undefined") {
      document.innerHTML = JSON.stringify({status: 400, data: "Bad argument(s)"});
    } else {
      const code = decodeURIComponent(query.code);

      const res = runGorilla(code);
      document.innerHTML = JSON.stringify({status: 200, data: res});
    }
  }
});