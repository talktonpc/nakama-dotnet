/**
* Copyright 2021 The Nakama Authors
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/

namespace NakamaSync
{
    /// <summary>
    /// A variable whose single value is synchronized across all clients connected to the same match.
    /// </summary>
    public class SelfVar<T> : Var<T>
    {
        public SelfVar(long opcode) : base(opcode)
        {
        }

        public void SetValue(T value)
        {
            // TODO need to handle sync match not started yet
            this.SetLocalValue(SyncMatch.Self, value);
        }

        internal override ISerializableVar<T> ToSerializable(bool isAck)
        {
            var serializable = base.ToSerializable(isAck);
            var presenceSerializable = new SerializablePresenceVar<T>(serializable, SyncMatch.Self.UserId);
            return presenceSerializable;
        }
    }
}
