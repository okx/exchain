exchaincli tx gov submit-proposal update-contract-bytecode ./contractblock.json --from captain --fees 0.001okt -y -b block --gas 8000000


sleep 8

exchaincli tx gov vote 1 yes --from captain --fees 0.001okt -y -b block --gas 8000000
