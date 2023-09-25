function gradingStudents(grades) {
    const roundedGrades = [];

    for (let i = 0; i < grades.length; i++) {
        const grade = grades[i];
        if (grade < 38) {
            roundedGrades.push(grade);
        } else {
            const nextMultipleOf5 = Math.ceil(grade / 5) * 5;
            if (nextMultipleOf5 - grade < 3) {
                roundedGrades.push(nextMultipleOf5);
            } else {
                roundedGrades.push(grade);
            }
        }
    }

    return roundedGrades;
}

function main() {
    // NOTES :
    // Run the code in terminal by script this "node main.js"
    
    const readline = require('readline');
    const rl = readline.createInterface({
        input: process.stdin,
        output: process.stdout
    });
    console.log('Input : ');
    rl.question('Enter the number of students: ', (n) => {
        const grades = [];
        let count = 0;

        function inputGrade() {
            rl.question(`Enter the grade for student ${count + 1}: `, (grade) => {
                grades.push(parseInt(grade));
                count++;

                if (count < n) {
                    inputGrade();
                } else {
                    rl.close();
                    const roundedGrades = gradingStudents(grades);
                    console.log('Output : ');
                    console.log(roundedGrades.join('\n'));
                }
            });
        }

        inputGrade();
    });
}

main();