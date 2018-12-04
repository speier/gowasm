if (!WebAssembly.instantiateStreaming) { // polyfill
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
        const source = await (await resp).arrayBuffer();
        return await WebAssembly.instantiate(source, importObject);
    };
}
(async () => {
    const go = new Go();
    const { instance } = await WebAssembly.instantiateStreaming(fetch('/app.wasm'), go.importObject);
    await go.run(instance);
})();