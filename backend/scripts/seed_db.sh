db=$1
dir=$2

files=$(ls "$dir" | grep sql | sort)

for f in $files 
do
    psql -d "${db}" -a -f "${dir}/${f}"
done;
