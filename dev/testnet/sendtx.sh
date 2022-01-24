# send tx

# for some client to broadcast txs 
# exchaincli tx send --node tcp://localhost:26657 but only stdtx 

MULTIVERSION=false
VERSIONS=("dev", "dev", "dev", "dev")
while getopts "sV:" opt; do 
    case $opt in
        V)
            echo "MULTIVERSION"
            if [[ -z $OPTARG ]]; then 
                MULTIVERSION=false 
            else 
                MULTIVERSION=true
                IFS=','
                read -ra VERSIONS <<< "$OPTARG"
            fi 
            ;;
    esac 
done

echo $MULTIVERSION

for val in "${VERSIONS[@]}";
do
  printf "name = $val\n"
done