echo "Starting Elevator script $0"
display_usage() {
  echo
  echo "Usage: $0"
  echo
  echo "Enter the port number - 5 digits"
  echo "Example ./safeStartWithParsing.sh 10001"
  echo
}
raise_error() {
  local error_message="$@"
  echo "${error_message}" 1>&2;
}
port="$1"
if [[ -z $port ]] ; then
  raise_error "Expected argument to be present"
  display_usage
  else 
    while :
    do
    echo "Starting the elevator"
    go run main.go -port=$port
    done
fi

