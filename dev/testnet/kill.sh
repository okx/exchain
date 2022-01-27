function killbyname_gracefully() {
  NAME=exchaind
  ps -ef|grep "$NAME"|grep -v grep |awk '{print "kill  "$2", "$8}'
  ps -ef|grep "$NAME"|grep -v grep |awk '{print "kill  "$2}' | sh
  echo "All <$NAME> killed gracefully!"
}
killbyname_gracefully