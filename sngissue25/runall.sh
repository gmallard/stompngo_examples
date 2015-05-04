#!/bin/bash
#
set -x
for d in client01-hang client02-hang client03-hang;do
	cd $d
	./run.sh
	cd ..
done
#
for d in clientindividual01-hang clientindividual02-hang clientindividual03-hang;do
	cd $d
	./run.sh
	cd ..
done
set +x

