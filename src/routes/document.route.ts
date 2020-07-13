import express from 'express'

const router = express.Router()

router.get('/', (req, res) => {
  res.send({ message: 'body' })
})

export const prefix = 'document'
export default router
