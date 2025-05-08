import { Request, Response } from 'express';
import { getRepository } from 'typeorm';
import { z } from 'zod';
import { Invoice } from './models';
import Mongo_Invoice from './mongoose';

// Zod schema for invoice validation
export const InvoiceSchema = z.object({
	vendor: z.string().min(1, 'Vendor name is required').max(255),
	amount: z.number().positive('Amount must be positive'),
	dueDate: z.string().refine(
		(value) => !isNaN(Date.parse(value)),
		{ message: 'Due date must be a valid date string' }
	),
	fileName: z.string().min(1, 'File name is required').max(255),
	rawInvoice: z.record(z.any()).optional(),
	isPaid: z.boolean().optional().default(false),
	paidDate: z.string().optional().nullable()
		.refine(
			(value) => value === null || value === undefined || !isNaN(Date.parse(value)),
			{ message: 'Paid date must be a valid date string or null' }
		)
});


export const addInvoice = async (req: Request, res: Response): Promise<void> => {
	try {
		// Validate request body against the schema
		const validatedData = InvoiceSchema.parse(req.body);

		// Create a new invoice entity
		const invoice = new Invoice()

		invoice.vendor = validatedData.vendor
		invoice.amount = validatedData.amount
		invoice.dueDate = new Date(validatedData.dueDate)
		invoice.fileName = validatedData.fileName
		invoice.rawInvoice = validatedData.rawInvoice || {}
		invoice.isPaid = validatedData.isPaid
		//invoice.paidDate = new Date(validatedData.paidDate || "")

		const newInvoice = Invoice.create(invoice);
		await newInvoice.save()

		//mongoose save
		const mongoInvoiceData = {
			vendor: newInvoice.vendor,
			amount: newInvoice.amount,
			dueDate: newInvoice.dueDate,
			fileName: newInvoice.fileName,
			rawInvoice: newInvoice.rawInvoice,
			isPaid: newInvoice.isPaid,
			externalId: newInvoice.id,
		};

		console.log(mongoInvoiceData)

		const mongo_invoice = new Mongo_Invoice(mongoInvoiceData)
		const doc = await mongo_invoice.save()
		console.log(doc._id)

		// Return the created invoice with a 201 status code
		res.status(201).json({
			success: true,
			message: 'Invoice created successfully',
			data: newInvoice,
		});
	} catch (error) {
		// Handle validation errors
		if (error instanceof z.ZodError) {
			res.status(400).json({
				success: false,
				message: 'Validation failed',
				errors: error.errors
			});
		}

		// Handle other errors
		console.error('Error creating invoice:', error);
		res.status(500).json({
			success: false,
			message: 'Failed to create invoice'
		});
	}

}


export const getInvoice = async (req: Request, res: Response): Promise<void> => {

	const { id } = req.params;

	try {

		// Validate that ID is provided
		if (!id) {
			res.status(400).json({
				success: false,
				message: 'Invoice ID is required'
			});
		}


		// Find invoice by ID
		const invoice = await Invoice.findOne({
			where: { id }
		});

		// Check if invoice exists
		if (!invoice) {
			res.status(404).json({
				success: false,
				message: `Invoice with ID ${id} not found`
			});
		}

		const document = await Mongo_Invoice.find({ externalId: id })

		// Return the invoice
		res.status(200).json({
			success: true,
			data: invoice,
			document: document
		});
	} catch (error) {
		console.error('Error fetching invoice:', error);

		// Handle invalid UUID format
		if (error instanceof Error && error.message.includes('invalid input syntax')) {
			res.status(400).json({
				success: false,
				message: 'Invalid invoice ID format'
			});
		}

		// Handle other errors
		res.status(500).json({
			success: false,
			message: 'Failed to fetch invoice'
		});
	}
}

