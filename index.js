 //literal characters
// const hello = "hello world"
// let regex = /hello/

// console.log(regex.test(hello))

//character classess 

// const character = "hello world"
// let regex1 = /[aeiou]/g;

// console.log(character.match(regex1))


// . will match any character but a new line

// const str = "hallo, hello";
// const regex = /h.llo/g;
// console.log(str.match(regex)); // ['hallo', 'hello']


const str = "ho, hoo, hooo";
const regex = /ho*/g;
console.log(str.match(regex)); // ['ho', 'hoo', 'hooo']




