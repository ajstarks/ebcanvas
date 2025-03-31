#!/bin/sh
for i in $(cat cl)
do
	cd $i
	case $i in
		echart)
		./allcharts
		;;
		elections)
		./allelections
		;;
		*)
		./$i &
		;;
	esac
	cd ..
done
