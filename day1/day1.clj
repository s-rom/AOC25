(require '[clojure.java.io :as io])
(require '[clojure.string :as str])
(require 'clojure.main)

(defn read-lines [path]
  (with-open [r (io/reader path)]
    (doall (line-seq r))))

(defn parse-rotation [rotation-str]
    (let [
        dir (if (= "L" (subs rotation-str 0 1))  -1 1)
        num (Integer/parseInt (subs rotation-str 1))
    ]
    (* dir num)))


(defn sign [x]
    (cond 
        (neg? x) -1 
        (pos? x) 1
        :else 0
))

(defn rotate [current input]
    (mod (+ current input) 100))


(defn count-turns [current input]
    "Bruteforced solution to count the number of times the dial crosses 0"
    (let [
        s (sign input)
        rot (abs input)
    ]

    (loop [turns 0 rotation rot c current]
        (if (= 0 rotation)
            turns ;; return number of turns

            (let [
                dial (rotate c s) ;; next dial position
            ]

            (recur 
                ;; inc turns if dial is 0
                (if (= 0 dial) (inc turns) turns)
                (dec rotation)
                dial
            )
            )
        )
    )
    )
)

(defn count-dial-at-0 [data]
    (->> data
        (map parse-rotation)
        (reductions rotate 50)
        (filter #(= % 0))
        (count)
    )
)


(println "Part 1: " 
    (count-dial-at-0 (read-lines "input1.txt"))
)

(println "Part 2: "
    (let [
        data (read-lines "input1.txt")
        parsed-data (map parse-rotation data)
        rotations  (conj (vec parsed-data) 0)

        dial      ( ->> parsed-data
                    (reductions rotate 50)
                    (vec))
        ]

        (reduce + (map count-turns dial rotations))        
    )
)
