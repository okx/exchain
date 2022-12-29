exchaincli tx gov submit-proposal update-contract-bytecode ./contractblock.json --from captain --fees 0.001okt -y -b block


sleep 20

exchaincli tx gov vote 1 yes --from captain --fees 0.001okt -y -b block
