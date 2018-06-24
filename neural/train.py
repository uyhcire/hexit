import random

import tensorflow as tf
import training_game_pb2


RESIDUAL_FILTERS = 16
RESIDUAL_BLOCKS = 2
LEARNING_RATE = 0.1


def get_data(x, y_, z_):
    with open('/go/src/github.com/uyhcire/hexit/training_games/training_game.dat', 'rb') as f:
        data = f.read()
    training_game = training_game_pb2.TrainingGame()
    training_game.ParseFromString(data)
    move_snapshot = random.choice(training_game.moveSnapshots)

    # Batch size = 1
    return {
        x: [list(move_snapshot.squaresOccupiedByMyself) + list(move_snapshot.squaresOccupiedByOtherPlayer)],
        y_: [move_snapshot.normalizedVisitCounts],
        z_: [[+1 if move_snapshot.winner == training_game_pb2.TrainingGame.MYSELF else -1]],
    }


def main():
    session = tf.Session()

    # 25 inputs for the player to move, 25 for the other player
    x = tf.placeholder(tf.float32, [None, 5*5*2], name='boardInput')

    # Policy
    W_y = tf.Variable(tf.zeros([5*5*2, 5*5]), name='policyWeights')
    b_y = tf.Variable(tf.zeros([5*5]), name='policyBiases')
    y = tf.add(tf.matmul(x, W_y), b_y, name='policyLogits')

    # Value
    W_z = tf.Variable(tf.zeros([5*5*2, 1]), name='valueWeights')
    b_z = tf.Variable(tf.zeros([1]), name='valueBiases')
    z = tf.nn.tanh(tf.add(tf.matmul(x, W_z), b_z), name='valueHeadOutput')
    
    y_ = tf.placeholder(tf.float32, [None, 5*5])
    z_ = tf.placeholder(tf.float32, [None, 1])

    cross_entropy = tf.nn.softmax_cross_entropy_with_logits(labels=y_, logits=y)
    policy_loss = tf.reduce_mean(cross_entropy)
    mse_loss = tf.reduce_mean(tf.squared_difference(z_, z))

    regularizer = tf.contrib.layers.l2_regularizer(scale=0.0001)
    reg_term = tf.contrib.layers.apply_regularization(regularizer, [W_y, W_z])

    loss = policy_loss + mse_loss + reg_term

    opt_op = tf.train.MomentumOptimizer(learning_rate=LEARNING_RATE, momentum=0.9, use_nesterov=True)
    train_op = opt_op.minimize(loss)

    init = tf.global_variables_initializer()
    session.run(init)

    # One training step
    feed_dict = get_data(x, y_, z_)
    session.run([train_op, loss], feed_dict=feed_dict)


if __name__ == '__main__':
    main()
