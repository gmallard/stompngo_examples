
# CP will be the Java CLASSPATH
CP=.

# Main Artemis JAR files.
for j in $(ls -1 ./jars/*.jar);do
  # echo $j
  CP="${CP}:$j"
done
echo "CLASSPATH: ${CP}"
