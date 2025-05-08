FROM oven/bun:1.0 as base

# Set working directory
WORKDIR /app

# Copy package files
COPY package.json bun.lockb* ./

# Install dependencies with Bun
RUN bun install --frozen-lockfile

# Copy application code
COPY . .

# Build TypeScript code if needed
RUN bun run build

# Expose the port the app runs on
EXPOSE 3000

# Command to run the application
CMD ["bun", "run", "start"]
