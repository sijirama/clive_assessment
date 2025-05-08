import { Entity, Column, PrimaryGeneratedColumn, CreateDateColumn, UpdateDateColumn, BaseEntity } from 'typeorm';

@Entity('invoices')
export class Invoice extends BaseEntity {
	@PrimaryGeneratedColumn('uuid')
	id!: string;

	@Column({ length: 255 })
	vendor!: string;

	@Column('decimal', { precision: 10, scale: 2 })
	amount!: number;

	@Column({ type: 'date' })
	dueDate!: Date;

	@Column({ length: 255, name: 'file_name' })
	fileName!: string;

	@Column('json', { nullable: true })
	rawInvoice!: Record<string, any>;

	@Column({ default: false })
	isPaid!: boolean;

	@Column({ type: 'date', nullable: true })
	paidDate!: Date | null;

	@CreateDateColumn({ name: 'created_at' })
	createdAt!: Date;

	@UpdateDateColumn({ name: 'updated_at' })
	updatedAt!: Date;
}


