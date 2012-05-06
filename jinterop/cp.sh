
# ActiveMQ lib directory
AHL=/ad3/amq/lib

# CP will be the Java CLASSPATH
CP=.

# Main AMQ JAR files
for j in $(ls -1 $AHL/*.jar);do
  # echo $j
  CP="${CP}:$j"
done

# Need optional JARs as well
for j in $(ls -1 $AHL/optional/*.jar);do
  # echo $j
  CP="${CP}:$j"
done

