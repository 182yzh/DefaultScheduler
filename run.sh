main() {

    cd $HOME

    if [ ! -d  log/dslog ];then
    mkdir log/dslog
    else
    echo dir "log/dslog" exist
    fi

    v=101
    while getopts "v:" opt; do
        if [[ $opt == "v" ]];then
            v=$OPTARG
            #echo "optarg is " ${OPTARG}
        fi
    done

    ### run the program
    cd -
    $HOME/DefaultScheduler/defaultscheduler -v=$v -log_dir=$HOME/log/dslog  -kubeconfig="/home/yzh/.kube/config" -totalgpu=1280 -gpulimit=64
    #./SchedulerFrame -v=$v -log_dir=$HOME/log  -kubeconfig="/home/yzh/.kube/config"   #-namespace="yzh" -alsologtostderr -gpulimit=1-go
}

main "$@"