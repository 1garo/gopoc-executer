const fs = require('fs');

// fs.writeFile("/tmp/test", "Hey there!", function(err) {
//     if(err) {
//         return err;
//     }
//     console.log("The file was saved!");
// }); 

// Or
fs.writeFileSync('/tmp/test-sync', 'Hey there!');
