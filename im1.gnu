set terminal 'png' font 'TimesNewRoman, 15'
set output "output.png"

set key above width 1 height 1

set xtics rotate by 45 right offset 0
set xtics time format "%tH:%tM:%tS"
set xrange [0:10000]
set yrange [0:120]

set grid xtics ytics mytics
set xlabel "Time"
set ylabel "Image Pulling Time(ms)"

set style line 1 lc rgb "#ef597b" pt 6 ps 1
set style line 2 lc rgb "#ff6d31" pt 5 ps 1
set style line 3 lc rgb "#73b66b" pt 4 ps 1
set style line 4 lc rgb "#ffcb18" pt 3 ps 1
set style line 5 lc rgb "#29a2c6" pt 2 ps 1

plot  "result5.txt" using 1:2 with points title "edge5" ls 5, \
      "result4.txt" using 1:2 with points title "edge4" ls 4, \
      "result3.txt" using 1:2 with points title "edge3" ls 3, \
      "result2.txt" using 1:2 with points title "edge2" ls 2, \
      "result1.txt" using 1:2 with points title "edge1" ls 1
