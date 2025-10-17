#make a longer shell script that computes the factorial of a number
#!/bin/bash

factorial() {
    if [ $1 -le 1 ]; then
        echo 1
    else
        local prev=$(factorial $(( $1 - 1 )))
        echo $(( $1 * prev ))
    fi
}

read -p "Enter a number: " num
result=$(factorial $num)
echo "The factorial of $num is $result"