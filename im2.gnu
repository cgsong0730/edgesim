set terminal 'png' font 'TimesNewRoman, 15'
set output "im2.png"

set boxwidth 1
set xlabel "Number of edge servers"
set ylabel "Image pulling time(ms)"
set key above width 1 height 1

set style fill solid 1.0 border lt -1

set style histogram gap 2
set style data histograms
set yrange [0:100]

plot "data_2_1.txt" using 2:xticlabels(1) lc rgb "#ef597b" title "non-mkrp", \
     "data_2_2.txt" using 2:xticlabels(1) lc rgb "#29a2c6" title "mkrp no", \
     "data_2_4.txt" using 2:xticlabels(1) lc rgb "#E69F00" title "comm", \
     "data_2_3.txt" using 2:xticlabels(1) lc rgb "#73b66b" title "mkrp no+af"
