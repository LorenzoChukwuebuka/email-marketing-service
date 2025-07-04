// Random Name Generator
export function generateRandomName() {
  const firstNames:string[] = [
    'Alex', 'Jordan', 'Taylor', 'Morgan', 'Casey', 'Riley', 'Avery', 'Quinn',
    'Sam', 'Blake', 'Parker', 'Reese', 'Sage', 'River', 'Skylar', 'Rowan',
    'Kai', 'Phoenix', 'Ember', 'Nova', 'Luna', 'Aria', 'Zoe', 'Maya',
    'Ethan', 'Liam', 'Noah', 'Oliver', 'James', 'Lucas', 'Mason', 'Logan',
    'Emma', 'Olivia', 'Ava', 'Isabella', 'Sophia', 'Mia', 'Charlotte', 'Amelia',
    'Harper', 'Evelyn', 'Abigail', 'Emily', 'Ella', 'Elizabeth', 'Camila', 'Luna',
    'Sofia', 'Avery', 'Mila', 'Aria', 'Scarlett', 'Penelope', 'Layla', 'Chloe',
    'Victoria', 'Madison', 'Eleanor', 'Grace', 'Nora', 'Riley', 'Zoey', 'Hannah',
    'Hazel', 'Lily', 'Ellie', 'Violet', 'Lillian', 'Zoe', 'Stella', 'Aurora',
    'Natalie', 'Emilia', 'Everly', 'Leah', 'Aubrey', 'Willow', 'Addison', 'Lucy',
    'Audrey', 'Bella', 'Nova', 'Brooklyn', 'Paisley', 'Savannah', 'Claire', 'Skylar',
    'Isla', 'Genesis', 'Naomi', 'Elena', 'Caroline', 'Eliana', 'Anna', 'Maya',
    'Valentina', 'Ruby', 'Kennedy', 'Ivy', 'Ariana', 'Aaliyah', 'Cora', 'Madelyn',
    'Alice', 'Kinsley', 'Hailey', 'Gabriella', 'Allison', 'Gianna', 'Serenity', 'Samantha',
    'Sarah', 'Autumn', 'Quinn', 'Eva', 'Piper', 'Sophie', 'Sadie', 'Delilah',
    'Josephine', 'Nevaeh', 'Adeline', 'Arya', 'Emery', 'Lydia', 'Clara', 'Vivian',
    'Madeline', 'Peyton', 'Julia', 'Rylee', 'Brielle', 'Reagan', 'Natalia', 'Jade',
    'Athena', 'Maria', 'Leilani', 'Everleigh', 'Liliana', 'Melanie', 'Mackenzie', 'Hadley',
    'Aiden', 'Jackson', 'Carter', 'Wyatt', 'Jayden', 'Dylan', 'Grayson', 'Levi',
    'Isaac', 'Gabriel', 'Julian', 'Mateo', 'Anthony', 'Jaxon', 'Lincoln', 'Joshua',
    'Christopher', 'Andrew', 'Theodore', 'Caleb', 'Ryan', 'Asher', 'Nathan', 'Thomas',
    'Leo', 'Isaiah', 'Charles', 'Josiah', 'Sebastian', 'Henry', 'Aaron', 'Eli',
    'Landon', 'Adrian', 'Cayden', 'Jordan', 'Gavin', 'Carson', 'Jace', 'Abbas',
    'Cameron', 'Connor', 'Santiago', 'Greyson', 'Jason', 'Ian', 'Nolan', 'Hunter',
    'Dominic', 'Cooper', 'Jude', 'Colton', 'Easton', 'Ezra', 'Jameson', 'Axel',
    'Zion', 'Rowan', 'Kai', 'Felix', 'Beckett', 'Damian', 'Maximus', 'Silas'
  ];

  const lastNames:string[] = [
    'Smith', 'Johnson', 'Williams', 'Brown', 'Jones', 'Garcia', 'Miller', 'Davis',
    'Rodriguez', 'Martinez', 'Hernandez', 'Lopez', 'Gonzalez', 'Wilson', 'Anderson', 'Thomas',
    'Taylor', 'Moore', 'Jackson', 'Martin', 'Lee', 'Perez', 'Thompson', 'White',
    'Harris', 'Sanchez', 'Clark', 'Ramirez', 'Lewis', 'Robinson', 'Walker', 'Young',
    'Allen', 'King', 'Wright', 'Scott', 'Torres', 'Nguyen', 'Hill', 'Flores',
    'Green', 'Adams', 'Nelson', 'Baker', 'Hall', 'Rivera', 'Campbell', 'Mitchell',
    'Carter', 'Roberts', 'Gomez', 'Phillips', 'Evans', 'Turner', 'Diaz', 'Parker',
    'Cruz', 'Edwards', 'Collins', 'Reyes', 'Stewart', 'Morris', 'Morales', 'Murphy',
    'Cook', 'Rogers', 'Gutierrez', 'Ortiz', 'Morgan', 'Cooper', 'Peterson', 'Bailey',
    'Reed', 'Kelly', 'Howard', 'Ramos', 'Kim', 'Cox', 'Ward', 'Richardson',
    'Watson', 'Brooks', 'Chavez', 'Wood', 'James', 'Bennett', 'Gray', 'Mendoza',
    'Ruiz', 'Hughes', 'Price', 'Alvarez', 'Castillo', 'Sanders', 'Patel', 'Myers',
    'Long', 'Ross', 'Foster', 'Jimenez', 'Powell', 'Jenkins', 'Perry', 'Russell',
    'Sullivan', 'Bell', 'Coleman', 'Butler', 'Henderson', 'Barnes', 'Gonzales', 'Fisher',
    'Vasquez', 'Simmons', 'Romero', 'Jordan', 'Patterson', 'Alexander', 'Hamilton', 'Graham',
    'Reynolds', 'Griffin', 'Wallace', 'Moreno', 'West', 'Cole', 'Hayes', 'Bryant'
  ];

  const nicknames:string[] = [
    'Ace', 'Dash', 'Zephyr', 'Raven', 'Storm', 'Blaze', 'Echo', 'Sage',
    'Frost', 'Viper', 'Ghost', 'Hawk', 'Wolf', 'Fox', 'Bear', 'Tiger',
    'Phoenix', 'Dragon', 'Falcon', 'Eagle', 'Spider', 'Venom', 'Shadow',
    'Bolt', 'Flash', 'Rocket', 'Bullet', 'Arrow', 'Blade', 'Steel',
    'Diamond', 'Ruby', 'Jade', 'Onyx', 'Amber', 'Crystal', 'Pearl',
    'Smokey', 'Dusty', 'Rusty', 'Sunny', 'Cloudy', 'Rainy', 'Misty',
    'Spike', 'Buzz', 'Fizz', 'Pop', 'Snap', 'Crackle', 'Boom',
    'Whisper', 'Shout', 'Hush', 'Roar', 'Growl', 'Purr', 'Chirp',
    'Doodle', 'Sketch', 'Paint', 'Brush', 'Canvas', 'Pixel', 'Code',
    'Byte', 'Chip', 'Circuit', 'Wire', 'Spark', 'Glow', 'Shine'
  ];

  // Random selection method
  const methods = ['fullName', 'nickname', 'firstNameOnly'];
  const method = methods[Math.floor(Math.random() * methods.length)];

  switch (method) {
    case 'fullName':
      const firstName = firstNames[Math.floor(Math.random() * firstNames.length)];
      const lastName = lastNames[Math.floor(Math.random() * lastNames.length)];
      return `${firstName} ${lastName}`;
    
    case 'nickname':
      return nicknames[Math.floor(Math.random() * nicknames.length)];
    
    case 'firstNameOnly':
      return firstNames[Math.floor(Math.random() * firstNames.length)];
    
    default:
      return firstNames[Math.floor(Math.random() * firstNames.length)];
  }
}

// Alternative function that generates multiple names
function generateMultipleNames(count = 10) {
  const names:string[] = [];
  for (let i = 0; i < count; i++) {
    names.push(generateRandomName());
  }
  return names;
}

// Usage examples:
console.log("Single random name:", generateRandomName());
console.log("10 random names:", generateMultipleNames(10));
console.log("100 random names:", generateMultipleNames(100));

// Export for use in other files
// module.exports = { generateRandomName, generateMultipleNames };