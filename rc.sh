#!/bin/sh
for i in $(cat cl)
do
	cd $i
	case $i in
		ebdeck)
		./ebdeck test.xml &
		;;
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
