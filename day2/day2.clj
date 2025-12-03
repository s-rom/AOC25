
(require '[clojure.java.io :as io])
(require '[clojure.string :as str])
(require 'clojure.main)

(defn read-lines [path]
  (with-open [r (io/reader path)]
    (doall (line-seq r))))

(defn invalid-id [id pattern]
  (not (= nil (re-find (re-pattern pattern) (str id)))))

(defn parse-data [raw]
  (map #(str/split % #"-") (str/split raw #",")))


(defn invalid-ids-in-range [[start end] pattern]

  (let [
    s (Long/parseLong start)
    e (Long/parseLong end)
  ]
    (reduce + (filter #(invalid-id % pattern) (range s (inc e))))))


(defn part-1 []
  (let [
    data (parse-data (first (read-lines "input.txt")))
  ]
   (reduce + (map #(invalid-ids-in-range % #"^(\d+)\1$") data ))))


(defn part-2 []
  (let [
    data (parse-data (first (read-lines "input.txt")))
  ]
   (reduce + (map #(invalid-ids-in-range % #"^(\d+)\1+$") data ))))

(println "Part 1: " (part-1))
(println "Part 2: " (part-2))
