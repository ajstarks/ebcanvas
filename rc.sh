#!/bin/sh
for i in $(cat cl)
do
	cd $i
	case $i in
		echart)
		./allcharts
		;;
		elections)
		./elections nyt-????.d  &
		;;
		*)
		./$i &
		;;
	esac
	cd ..
done
