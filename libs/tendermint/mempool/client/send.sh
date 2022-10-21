./client --buy=true --amount=1 --priv=8ff3ca2d9985c3a52b459e2f6e7822b23e1af845961e22128d5f372fb9aa5f17
./client --buy=true --amount=1 --priv=e47a1fe74a7f9bfa44a362a3c6fbe96667242f62e6b8e138b3f61bd431c3215d
./client --buy=false --amount=3 --priv=75dee45fc7b2dd69ec22dc6a825a2d982aee4ca2edd42c53ced0912173c4a788
sleep 2

./client --buy=false --amount=5 --price=18200000000000000000001
./client --buy=false --amount=31 --price=18200000000000000000002
./client --buy=false --amount=212 --price=18200000000000000000002
sleep 2

./client --buy=false --amount=111 --price=18200000000000000000002
./client --buy=false --amount=10 --price=18200000000000000000002
./client --buy=false --amount=10 --price=18200000000000000000003
./client --buy=false --amount=1 --price=18200000000000000000003
sleep 2
./client --buy=false --amount=123 --price=18200000000000000000003
./client --buy=false --amount=9 --price=18200000000000000000004
./client --buy=false --amount=63 --price=18200000000000000000003
./client --buy=false --amount=108 --price=18200000000000000000003
sleep 2
./client --buy=false --amount=63 --price=18199999999999999999999
./client --buy=false --amount=5 --price=18199999999999999999998
./client --buy=false --amount=41 --price=18199999999999999999997
./client --buy=true --amount=19 --price=18199999999999999999997 --priv=e47a1fe74a7f9bfa44a362a3c6fbe96667242f62e6b8e138b3f61bd431c3215d
sleep 2
./client --buy=true --amount=124 --price=18199999999999999999998 --priv=75dee45fc7b2dd69ec22dc6a825a2d982aee4ca2edd42c53ced0912173c4a788
./client --buy=true --amount=12 --price=18199999999999999999999 --priv=75dee45fc7b2dd69ec22dc6a825a2d982aee4ca2edd42c53ced0912173c4a788
./client --buy=true --amount=92 --price=18200000000000000000000 --priv=e47a1fe74a7f9bfa44a362a3c6fbe96667242f62e6b8e138b3f61bd431c3215d
