##### store code  
exchaincli tx wasm store ./target/wasm32-unknown-unknown/release/callee.wasm --from captain --gas-prices 0.0000000001okt --gas auto -b block --gas-adjustment 1.5 -y  
##### instantiate  
exchaincli tx wasm instantiate 1 '{}' --from captain --gas-prices 0.0000000001okt --gas auto -b block --gas-adjustment 1.5 -y  
##### execute  
exchaincli tx wasm execute 0x5A8D648DEE57b2fc90D98DC17fa887159b69638b '{"add":{"delta":"16"}}'  --from captain --gas-prices 0.0000000001okt --gas auto -b block --gas-adjustment 1.5 -y  
##### query  
exchaincli query wasm contract-state smart 0x5A8D648DEE57b2fc90D98DC17fa887159b69638b '{"get_counter":{}}'  

