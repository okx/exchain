package types


// WasmConfig is the extra config required for wasm
type WasmConfig struct {}

// DefaultWasmConfig returns the default settings for WasmConfig
func DefaultWasmConfig() WasmConfig {
	return WasmConfig{}
}
