rm -rf s0
cp -R s0_bak s0
exchaind replay --home ./s0 -d ./sx/data
