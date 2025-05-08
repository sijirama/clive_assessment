import mongoose, { Schema, Document } from 'mongoose';

// Define the Invoice interface
export interface IInvoice extends Document {
	externalId: string;
	vendor: string;
	amount: number;
	dueDate: Date;
	fileName: string;
	rawInvoice?: Record<string, any>;
	isPaid: boolean;
	paidDate?: Date | null;
	createdAt: Date;
	updatedAt: Date;
}

// Create the Invoice schema
const invoiceSchema = new Schema<IInvoice>(
	{
		externalId: {
			type: String,
			required: true
		},
		vendor: {
			type: String,
			required: true,
			maxlength: 255
		},
		amount: {
			type: Number,
			required: true,
			min: 0,
			get: (v: number) => parseFloat(v.toFixed(2)), // Handle decimal precision
			set: (v: number) => parseFloat(v.toFixed(2))
		},
		dueDate: {
			type: Date,
			required: true
		},
		fileName: {
			type: String,
			required: true,
			maxlength: 255
		},
		rawInvoice: {
			type: Schema.Types.Mixed,
			default: null
		},
		isPaid: {
			type: Boolean,
			default: false
		},
		paidDate: {
			type: Date,
			default: null
		}
	},
	{
		timestamps: {
			createdAt: 'createdAt',
			updatedAt: 'updatedAt'
		},
		toJSON: {
			getters: true,
			virtuals: true,
			transform: (_, ret) => {
				delete ret._id;
				return ret;
			}
		},
		id: true
	}
);

// Ensure decimal precision for amount
invoiceSchema.path('amount').get(function (num: number) {
	return parseFloat(num.toFixed(2));
});

// Create and export the model
const Invoice = mongoose.model<IInvoice>('Invoice', invoiceSchema);

export default Invoice;
