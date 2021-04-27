#!/bin/bash
# name gpuneed tasknums
#printf "%s %d %d\n" $1 $2 $3  

cp Test/gpu_spin.yaml  yaml/$1.yaml
# Task Name
sed -i '4s/replace-value/'$1'/'  yaml/$1.yaml
sed -i '11s/replace-value/'$1'/'  yaml/$1.yaml
sed -i '20s/replace-value/'$1'/'  yaml/$1.yaml


# Task Resource Need
sed -i '27s/replace-value/'$2'/'  yaml/$1.yaml
sed -i '31s/replace-value/'$2'/'  yaml/$1.yaml

# Task Number
sed -i '6s/replace-value/'$3'/'  yaml/$1.yaml
sed -i '7s/replace-value/'$3'/'  yaml/$1.yaml
sed -i '16s/replace-value/"'$3'"/'  yaml/$1.yaml

#viour Task Time
sed -i '32s/replace-value/'$4'/'  yaml/$1.yaml


