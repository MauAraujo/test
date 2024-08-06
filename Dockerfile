# Step 1: Set the base image. Here, we're using the official Node.js 16 Alpine image as it's lightweight.
FROM node:18


# Step 2: Set the working directory inside the container
WORKDIR /app

# Step 3: Copy package.json and package-lock.json (or yarn.lock) files
COPY package*.json ./


# Step 4: Install dependencies
RUN npm install

# Step 5: Copy the rest of your app's source code
# COPY src .
COPY . /app

# Step 6: Compile TypeScript to JavaScript.
# Note: Your package.json must have a "build" script that compiles TypeScript.
RUN npm run build

# Step 7: Expose the port your app runs on
EXPOSE 3000

# Step 8: Define the command to run your app using the compiled JavaScript
# Note: Adjust "dist/index.js" based on your output directory and main file name
CMD [ "node", "dist/index.js" ]