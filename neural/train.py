import os
import random

import tensorflow as tf
import training_game_pb2


RESIDUAL_FILTERS = 16
RESIDUAL_BLOCKS = 2
LEARNING_RATE = 0.1


def get_data():
    inputs = []
    policy_targets = []
    value_targets = []

    i = 0
    while os.path.exists('training_games/{}'.format(i)):
        print('{} games loaded'.format(i))
        with open('training_games/{}'.format(i), 'rb') as f:
            data = f.read()
        training_game = training_game_pb2.TrainingGame()
        training_game.ParseFromString(data)
        move_snapshot = random.choice(training_game.moveSnapshots)
        inputs.append(
            list(move_snapshot.squaresOccupiedByMyself) + \
            list(move_snapshot.squaresOccupiedByOtherPlayer))
        policy_targets.append(list(move_snapshot.normalizedVisitCounts))
        value_targets.append([+1] if move_snapshot.winner == training_game_pb2.TrainingGame.MYSELF else [-1])
        i += 1

    return inputs, policy_targets, value_targets


def main():
    session = tf.Session()
    tf.keras.backend.set_session(session)

    # 25 inputs for the player to move, 25 for the other player
    board_input = tf.keras.layers.Input(shape=(5*5*2,), dtype='float32', name='boardInput')
    policy_output = tf.keras.layers.Dense(5*5, activation='softmax', name='policyOutput')(board_input)
    value_output = tf.keras.layers.Dense(1, activation='tanh', name='valueOutput')(board_input)

    model = tf.keras.models.Model(inputs=[board_input], outputs=[policy_output, value_output])
    sgd = tf.keras.optimizers.SGD(lr=LEARNING_RATE, momentum=0.9, nesterov=True)
    model.compile(
        optimizer=sgd, 
        loss=['categorical_crossentropy', 'mean_squared_error'],
        loss_weights=[1.0, 1.0],
    )

    print('Loading data...')
    inputs, policy_targets, value_targets = get_data()
    print('...done')

    model.fit(
        [inputs], [policy_targets, value_targets], 
        epochs=10, batch_size=100, validation_split=0.1)


if __name__ == '__main__':
    main()
