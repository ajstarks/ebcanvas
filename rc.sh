#!/bin/sh
for i in $(cat cl)
do
	cd $i
	case $i in
	       echart)
	           ./allcharts
				;;
			*)
	           ./$i &
				;;
	esac
	cd ..
done
