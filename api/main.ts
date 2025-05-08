import express, { Router } from 'express';
import { addInvoice, getInvoice } from './controllers';
import { DataSource } from 'typeorm';
import { Invoice } from './models'
import mongoose from 'mongoose';

export const AppDataSource = new DataSource({
	type: 'postgres',
	host: process.env.PG_HOST || 'postgres',
	port: parseInt(process.env.PG_PORT || '5432'),
	username: process.env.PG_USER || 'postgres',
	password: process.env.PG_PASSWORD || 'postgres',
	database: process.env.PG_DB || 'insight',
	entities: [Invoice],
	synchronize: true,
	logging: true,
});

const app = express();

app.use(express.json());
app.use(express.urlencoded({ extended: true }))

const router = Router();

router.post('/', addInvoice);
router.get('/:id', getInvoice)

app.use("/api/", router)


AppDataSource.initialize()
	.then(() => {
		console.log('Data Source has been initialized!');
		mongoose.connect("mongodb://mongo:27017/test").then(() => {

			console.log('Mongoose has been initialized!');


			app.listen(3000, () => {
				console.log(`Server running on port 3000`);
			});

		}).catch((error) => {
			console.error('Error during Mongoose initialization:', error);
			process.exit(1);

		})

	})
	.catch((error) => {
		console.error('Error during Data Source initialization:', error);
		process.exit(1);
	});
