import { Entity, Column, PrimaryColumn, CreateDateColumn } from 'typeorm'

@Entity()
export class Document {
  @PrimaryColumn()
  id: string

  @Column()
  content: string

  @CreateDateColumn()
  dateCreated: Date

  @Column({
    default: 'text'
  })
  extension: string
}
