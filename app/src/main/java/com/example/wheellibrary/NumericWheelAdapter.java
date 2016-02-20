/*
 *  Copyright 2010 Yuri Kanivets
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package com.example.wheellibrary;

/**
 * Numeric Wheel adapter.
 */
public class NumericWheelAdapter implements WheelAdapter {

    private int[] ops;

    /**
     * Default constructor
     */
    public NumericWheelAdapter(int[] ops) {
        this.ops = ops;
    }

    @Override
    public String getItem(int index) {
        return Integer.toString(ops[index]);
    }

    @Override
    public int getItemsCount() {
        return ops.length;
    }

    @Override
    public int getMaximumLength() {
        int max = -1, temp = -1;
        for (int op : ops) {
            temp = Integer.toString(op).length();
            if (temp > max) {
                max = temp;
            }
        }
        return max;
    }
}
