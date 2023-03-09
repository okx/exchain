#!/bin/bash

    for (( i=0; ; ))
    do
        okbchaincli tx send captain 0x83D83497431C2D3FEab296a9fba4e5FaDD2f7eD0 1okb --fees 1okb -b block -y
    done

