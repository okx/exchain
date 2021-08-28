#!/bin/bash
     ############################
     ########## 裁剪 #############
CURDIR=`dirname $0`
HOME=$CURDIR/"_cache_evm"
HEIGHT=50 # 从开头直到高度为height的区块全部裁掉，注意不包括height本身
ENABLE_PRUNING=true # 默认为true, 设置enable_pruning为false即可跳过pruning执行阶段

    # prune 然后 compact, 为防止compact一次执行达不到效果，程序内部默认执行五次调用
    # 裁剪application state，
exchaind data prune-compact state \
    --home ${HOME} \
    --height ${HEIGHT} \
    --enable_pruning=${ENABLE_PRUNING}

    # blocks and states，注意有些state可能不被允许删除则强制保留。
exchaind data prune-compact block \
    --home ${HOME} \
    --height ${HEIGHT} \
    --enable_pruning=${ENABLE_PRUNING}
    
    #############################
    ########## 查询 ##############
    # application state, 输出为数组格式
exchaind data query state --home ${HOME}
    # blocks and states，输出为闭区间格式
exchaind data query block --home ${HOME}