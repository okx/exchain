#### 1. There are four contracts(contract_simple.rs, contract_funds, contract_simple_instantiate, contract_use_return) in floder of src, the user could modify the src/lib.rs import file for compile contract, the default contract is contract_simple.rs.

#### 2. The command.

(1) when test **contract_simple.rs**  
##### store code  
exchaincli tx wasm store ./target/wasm32-unknown-unknown/release/caller.wasm --from captain --gas-prices 0.0000000001okt --gas auto -b block --gas-adjustment 1.5 -y  
##### instantiate  
exchaincli tx wasm instantiate 2 '{"addr":""}' --from captain --gas-prices 0.0000000001okt --gas auto -b block --gas-adjustment 1.5 -y  
##### execute, the field of addr is the callee address  
exchaincli tx wasm execute $new_contract_address '{"call":{"delta":"9", "addr":"0x5A8D648DEE57b2fc90D98DC17fa887159b69638b"}}' --from captain --gas-prices 0.0000000001okt --gas auto -b block --gas-adjustment 1.5 -y  
exchaincli tx wasm execute $new_contract_address '{"delegate_call":{"delta":"9", "addr":"0x5A8D648DEE57b2fc90D98DC17fa887159b69638b"}}' --from captain --gas-prices 0.0000000001okt --gas auto -b block --gas-adjustment 1.5 -y  
##### query  
exchaincli query wasm contract-state smart $new_contract_address '{"get_counter":{}}'  
exchaincli query wasm contract-state smart 0x5A8D648DEE57b2fc90D98DC17fa887159b69638b '{"get_counter":{}}'  

(2) when test **contract_funds.rs**  
##### store code  
exchaincli tx wasm store ./target/wasm32-unknown-unknown/release/caller.wasm --from captain --gas-prices 0.0000000001okt --gas auto -b block --gas-adjustment 1.5 -y  
##### instantiate  
exchaincli tx wasm instantiate 3 '{"addr":""}' --from captain --gas-prices 0.0000000001okt --gas auto -b block --gas-adjustment 1.5 -y  
##### send coins to the contract address  
exchaincli tx send captain $new_contract_address 100okt --fees 1okt -b block -y  
// execute, the field of addr is the callee address  
exchaincli tx wasm execute $new_contract_address '{"call":{"delta":"9", "addr":"0x5A8D648DEE57b2fc90D98DC17fa887159b69638b"}}' --from captain --gas-prices 0.0000000001okt --gas auto -b block --gas-adjustment 1.5 -y  
##### query  
exchaincli query wasm contract-state smart $new_contract_address '{"get_counter":{}}'  
exchaincli query wasm contract-state smart 0x5A8D648DEE57b2fc90D98DC17fa887159b69638b '{"get_counter":{}}'  

(3) when test **contract_simple_instantiate.rs**  
##### store code  
exchaincli tx wasm store ./target/wasm32-unknown-unknown/release/caller.wasm --from captain --gas-prices 0.0000000001okt --gas auto -b block --gas-adjustment 1.5 -y  
##### instantiate  
exchaincli tx wasm instantiate 4 '{"addr":"0x5A8D648DEE57b2fc90D98DC17fa887159b69638b"}' --from captain --gas-prices 0.0000000001okt --gas auto -b block --gas-adjustment 1.5 -y  
##### query  
exchaincli query wasm contract-state smart $new_contract_address '{"get_counter":{}}'  
exchaincli query wasm contract-state smart 0x5A8D648DEE57b2fc90D98DC17fa887159b69638b '{"get_counter":{}}'  

(4) when test **contract_use_return.rs**  
##### store code  
exchaincli tx wasm store ./target/wasm32-unknown-unknown/release/caller.wasm --from captain --gas-prices 0.0000000001okt --gas auto -b block --gas-adjustment 1.5 -y  
##### instantiate  
exchaincli tx wasm instantiate 5 '{"addr":""}' --from captain --gas-prices 0.0000000001okt --gas auto -b block --gas-adjustment 1.5 -y  
##### execute, the field of addr is the callee address  
exchaincli tx wasm execute $new_contract_address '{"call":{"delta":"9", "addr":"0x5A8D648DEE57b2fc90D98DC17fa887159b69638b"}}' --from captain --gas-prices 0.0000000001okt --gas auto -b block --gas-adjustment 1.5 -y  
exchaincli tx wasm execute $new_contract_address '{"delegate_call":{"delta":"9", "addr":"0x5A8D648DEE57b2fc90D98DC17fa887159b69638b"}}' --from captain --gas-prices 0.0000000001okt --gas auto -b block --gas-adjustment 1.5 -y  
##### query  
exchaincli query wasm contract-state smart $new_contract_address '{"get_counter":{}}'  
exchaincli query wasm contract-state smart 0x5A8D648DEE57b2fc90D98DC17fa887159b69638b '{"get_counter":{}}'  