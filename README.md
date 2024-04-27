# Telegram Bot Development

This README file outlines the steps taken to develop a Telegram bot with basic filtering functionality using Go (Golang) and PostgreSQL.

## Development Process

1. **Understanding Telegram Bot API**: Started by familiarizing myself with the Telegram Bot API and its functionalities. Utilized the template provided by Telegram and studied the API documentation.

2. **Learning Go (Golang)**: Engaged in Go tutorials to understand the basics of the language and its features. Utilized various online resources and tutorials to grasp the fundamentals.

3. **Learning PostgreSQL**: Followed tutorials and documentation to learn about PostgreSQL, an open-source relational database management system. Understood how to interact with PostgreSQL using SQL queries and command line tools.

4. **Structuring the Bot**: Decided to take an object-oriented programming (OOP) approach to keep the code organized and maintainable. Created separate structs for the Telegram bot and the PostgreSQL database.

5. **Client-Database Connection**: Established the functionality for the client to connect to the database, allowing interactions within the network.

6. **Filtering Mechanism**: Developed a basic filtering mechanism where users are required to enter a filter word to proceed with interactions.

7. **Security Measures**: Ensured that sensitive information such as bot token and database password were hidden to maintain security. Implemented methods to securely handle authentication and access control.

8. **Dependency Management**: Removed unnecessary dependencies from the project to optimize performance and reduce complexity.

9. **Adding Bot Commands**: Implemented various bot commands such as `/help` to provide assistance, and `/stop` to gracefully close the database connection and stop the bot.

10. **Bot Token Management**: Integrated functionality to handle bot token revocation and regeneration by the user.

11. **Enhancements**: Continued to enhance the bot's functionality based on requirements and user feedback. Added features like `/filter` command to retrieve the filtering word.

## Usage

To use the Telegram bot:

1. **Obtain a Bot Token**: Obtain a bot token from the BotFather on Telegram. Replace `YOUR_BOT_TOKEN` in the source code with your own bot token.

2. **Provide SQL Connection String**: Provide the SQL connection string in the format specified in the source code. Replace `YOUR_SQL_CONNECTION_STRING` with your own SQL connection string.

3. **Run the Bot Application**: Run the bot application, either locally or within a Docker container.

4. **Interact with the Bot**: Interact with the bot using various commands as described in the bot's documentation.

**Important**: Ensure that you use your own bot token and SQL connection string for security and customization purposes.
